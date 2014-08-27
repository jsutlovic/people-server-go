package main

func main() {
	NewServer(&Config{"postgres", "user=vagrant dbname=people host=/var/run/postgresql sslmode=disable application_name=people-go", "0.0.0.0:3000"}).Serve()
}
