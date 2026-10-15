# Green Tiles

Yet another web tool to view your GitHub contributions. Inspired by [Green Wall](https://green-wall.vercel.app/) and [a list of others](https://twitter.com/wanglei001/status/1610496029652324359).

![green-tiles](https://user-images.githubusercontent.com/2587202/218168500-b11af51a-0067-496c-9e2a-7573a04e5a05.png)

## Tech Stacks
* Framework: [Vite](https://vitejs.dev/)
* UI Library: [daisyUI](https://daisyui.com/)
* Backend: AWS API Gateway + AWS Lambda (here's a [tutorial](https://old-panda.com/2020/03/02/lambda-api-gateway-note/) in Chinese)

## Project Setup

Use Node.js 22.22.2 or newer and pnpm 10.14.0. The backend requires Go 1.27.2 or newer.

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

### Coverage and CI

```sh
pnpm run test:coverage
cd backend
go test -race -covermode=atomic -coverprofile=coverage.out ./...
bash ../.github/scripts/check-go-coverage.sh coverage.out
```

The frontend requires at least **95% statements, branches, functions, and lines in
each source file**, including Vue components and the app entry point. Type
declarations and test files are excluded. The backend requires at least **95%
statement coverage both overall and in each package**. Tests stub GitHub requests
and do not require credentials or live API access.

GitHub Actions runs these coverage gates on pushes and pull requests, along with
type checks, production builds, dependency audits, Go vet, and Go race detection.
Coverage reports are available as workflow artifacts; the frontend HTML report is
also generated locally at `coverage/index.html`.
