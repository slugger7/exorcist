# Exorcist Web

Web frontend for the exorcist project

## Under construction

This is your warning that this project is not in any state ready to be used (even by me) and is still under constant construction without a usable version at the moment

## Getting started

### Prerequisites

- [Exorcist server](https://github.com/slugger7/exorcist)
- Node JS

### Run

- Run [Exorcist server](https://github.com/slugger7/exorcist) following their startup guide
- Copy [.env.example](./.env.example) to `.env` and modify parameters
- `npm i`
- `npm run dev`

## Vision

Having a completely independent frontend project from the backend while still having certainty that both projects will work with each other without depending on each other.

This should be achieveable using a mock server to stub out the actual backend. This gives you a lot of control in order to change things quickly to get the frontend workning while not actually having to implement it.

In order to make sure that the frontend and the backend can still communicate effectively, contract based testing will have to be implemented in order to ensure compatibility.
