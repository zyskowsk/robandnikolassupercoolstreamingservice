package main

// Interface for Server-like objects
type IRunnable interface {
	Run(port int32) (error)
	Close()

	Process_request(req SongRequest) (error)
}

type IGetSong interface {
	Get_song_chunk(req SongChunkRequest) (SongChunk, error)
}