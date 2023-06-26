go build -o bookings cmd/web/*.go
./bookings -dbname=bookings -cache=false -production=false
