echo "Dependancies..."

dep ensure -v

echo "Cleanup Previous Builds"
rm -rf ../bin/

echo "Building Image Generator..."
env GOOS=linux go build -ldflags="-s -w" -o ../bin/genimage ../src/genimage

echo "Building Image Generator.."
cd src/tamadb/handlers/
for f in *.go; do
  filename="${f%.go}"
  if GOOS=linux go build -o "../../../bin/tamadb/handlers/$filename" ${f}; then
    echo "✓ Compiled $filename"
  else
    echo "✕ Failed to compile $filename!"
    exit 1
  fi
done

echo "Done."