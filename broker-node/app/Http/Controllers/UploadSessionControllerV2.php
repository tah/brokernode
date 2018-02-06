<?php

namespace App\Http\Controllers;

use App\Clients\BrokerNode;
use App\DataMap;
use App\UploadSession;
use GuzzleHttp\Client;
use Illuminate\Http\Request;
use Illuminate\Support\Collection;
use Tuupola\Trytes;

class UploadSessionControllerV2 extends Controller
{
    public function update(Request $request, $id)
    {
        $session = UploadSession::find($id);
        if (empty($session)) return response('Session not found.', 404);

        $genesis_hash = $session['genesis_hash'];
        $chunks = $request->input('chunks');

        $res_addr = "{$_SERVER['REMOTE_ADDR']}:{$_SERVER['REMOTE_PORT']}";

        // Adapt chunks to reqs that hooknode expects.
        $chunk_reqs = collect($chunks)
          ->map(function($chunk, $idx) {
            // TODO: Adapt to whatever processChunks expects.
            $trytes = new Trytes(["characters" => Trytes::IOTA]);
            $hash_in_trytes = $trytes->encode($chunk["hash"]);
            $shortened_hash = substr($hash_in_trytes, 0, 81);

            return [
              'response_address' => $res_addr,
              'address' => $shortened_hash,
              'message' => $chunk['data'],
              'chunk_idx' => $chunk['idx'],
            ];
          });

        // Process chunks in 1000 chunk batches.
        $chunk_reqs
          ->chunk(1000) // Limited by IRI API.
          ->each(function($req_list, $idx) {
            BrokerNode::processChunks($req_list);
          });

        // Save to DB.
        $chunk_reqs
          ->each(function($req, $idx) {
            DataMap::where('genesis_hash', $genesis_hash)
              ->where('chunk_idx', $req['chunk_idx'])
              ->update([
                'address' => $req['address'],
                'message' => $req['message'],
              ]);
          });

        return response('Success.', 204);
    }
}