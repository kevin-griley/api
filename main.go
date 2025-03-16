package main

func main() {

	apiConfig := NewApiConfig()
	apiServer := NewAPIServer(apiConfig)
	apiServer.Run()

}
