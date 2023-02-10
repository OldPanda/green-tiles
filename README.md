# Green Tiles

Yet another web tool to view your GitHub contributions. Inspired by [Green Wall](https://green-wall.vercel.app/) and [a list of others](https://twitter.com/wanglei001/status/1610496029652324359).

![green-tiles](https://user-images.githubusercontent.com/2587202/218168500-b11af51a-0067-496c-9e2a-7573a04e5a05.png)

## Tech Stacks
* Framework: [Vite](https://vitejs.dev/)
* UI Library: [daisyUI](https://daisyui.com/)
* Backend: AWS API Gateway + AWS Lambda (here's a [tutorial](https://old-panda.com/2020/03/02/lambda-api-gateway-note/) in Chinese)

## Project Setup

```sh
pnpm install
```

### Compile and Hot-Reload for Development

```sh
pnpm run dev
```

### Type-Check, Compile and Minify for Production

```sh
pnpm run build
```

### Run Unit Tests with [Vitest](https://vitest.dev/)

```sh
pnpm run test:unit
```
