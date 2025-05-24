# order-packs-calculator

Find and optimal way to pack your gear!
A Go project that calculates the optimal packing of items into boxes of different sizes, minimizing waste and optimizing
efficiency.

## Features

- Calculates the most efficient way to pack items into boxes
- Minimizes waste by using the fewest boxes possible
- Handles large order sizes efficiently using dynamic programming
- Validates input parameters for order size and pack configurations

## I want to see the demo!

The repository is public. In order to access a publicly available demo, please reach out to the appropriate point of
contact.

## Running locally

   ```bash
   docker compose up
   ```  
Once the container is up and running, the demo page will be available at `http://localhost:8080`.

## Running the tests

   ```bash
   go test ./...
   ```

## Limitations

The application supports order sizes up to 1 million items.

