# Set the app to deploy to
export HEROKU_APP=ovrstat

# Remove the vendor folder to fetch fresh dependencies
rm -rf vendor

# Install all latest dependencies
glide cache-clear
glide install

# Build the binary that will run in the Docker container
env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/ovrstat

# Build and deploy the Docker image to Heroku
heroku container:login
heroku container:push web
heroku container:release web

# Delete the old binary
rm -rf ovrstat