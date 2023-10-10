<!-- markdownlint-disable-file MD013 -->
# The Million-ERC-20 DApp

A [Cartesi Rollups](https://github.com/cartesi/rollups) DApp inspired by the [Millon Dollar Homepage](http://www.milliondollarhomepage.com) and built with [EggRoll](https://github.com/gligneul/eggroll).

The idea is to provide a single-page application (SPA) that displays a `1000 px` x `1000 px` image that initially has all its pixels up for sale.

The SPA should display the current state of the image, as informed by the DApp, painting the pictures already owned by users at their corresponding places and all space on sale.

It should provide a means for any user to buy a chunk (a rectangle) of the image to display a PNG image linked to a give URL. 

Each pixel will cost 1 unit of a given preconfigured ERC-20 token, hence valuing the whole image in 1,000,000 units of the given ERC-20 token.

## Architecture

Being built with the help of eggroll, the DApp is going to be mostly written in Go and it's going to be comprised of the following components:

### `Contract`

The `Contract` defines the contract state an how to *advance* or *inspect* its partial state, which should keep control of the chunks of space still available and the ones already owned along, with the corresponding owner addresses.

> The state of the DApp kept in the `Contract` should be minimal, due to the limits defined by the Cartesi Machine.
> Hence the suggestion.
> Of course this may be changed during development.

### `Client`

The `Client` sends inputs to the `Contract` and keeps the full state of the DApp.

The state must keep the full picture, with all the chunks and their respective owners, along with the PNG and URL to be displayed.

The `Client` will also work as a web server and serve a SPA for user interaction.

#### `SPA`

The `SPA` will handle user submissions and will integrate with a wallet for handling payments.

> Most likely, the wallet integration is the only part of the DApp that will rely on a language other than Go.
