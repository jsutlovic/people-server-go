package main

func main() {
	NewServer(MustReadConfigFile("config.yml")).Serve()
}
