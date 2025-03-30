package main

func main() {
	tgServer := WireBuild()

	tgServer.Listen()
}
