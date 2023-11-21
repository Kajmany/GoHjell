# GoHjell: A Very Specific UCI Chess Engine Wrapper
*Unfortunately, the service doesn't yet work.*

## Required
- The desired UCI-compatible Chess engine, over PATH OR the ability to make programmatic requests to my hosted version.

## Getting Started
Currently, local port 8080 is hard-coded. Simply start the server and send HTTP POST requests with JSON bodies to the 
endpoint /analyze/! 

## Usage

A request might look like the following:

`{
"FEN": "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
"MultiPV": 3,
"Depth": 20
}`

The Forsyth-Edwards Notation must be valid non-extended format. MultiPV is configurable from 1-10 variations. Depth is configurable from 1-20.

When they're functioning, responses will look like:

`{
"variations": [
{
"depth": 20,
"score": 36,
"rank": 1,
"moves": ["d2d4", "d7d5", "c2c4", "e7e6", "g1f3", "g8f6", "g2g3", "d5c4", "d1a4", "c8d7", "a4c4", "d7c6", "f1g2", "c6d5", "c4d3", "d5e4", "d3b3", "b8c6", "e1g1", "c6d4", "f3d4", "d8d4"]
},
{
"depth": 20,
"score": 34,
"rank": 2,
"moves": ["g1f3", "d7d5", "d2d4", "g8f6", "c2c4", "c7c6", "b1c3", "d5c4", "a2a4", "e7e6", "e2e3", "c6c5", "f1c4", "c5d4", "e3d4", "f8e7", "d1e2"]
},
{
"depth": 20,
"score": 31,
"rank": 3,
"moves": ["e2e4", "e7e5", "g1f3", "b8c6", "f1b5", "a7a6", "b5a4", "g8f6", "e1g1", "f8e7", "f1e1", "b7b5", "a4b3", "e8g8", "h2h3", "c8b7", "d2d3", "d7d5", "e4d5", "f6d5"]
}
]
}`

My hosted version will require basic authentication, so don't forget to add the header with a supplied key!