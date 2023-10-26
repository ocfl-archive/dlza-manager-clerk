# create-svelte

Everything you need to build a Svelte project, powered by [`create-svelte`](https://github.com/sveltejs/kit/tree/master/packages/create-svelte).

## Creating a project

If you're seeing this, you've probably already done this step. Congrats!

```bash
# create a new project in the current directory
npm create svelte@latest

# create a new project in my-app
npm create svelte@latest my-app
```

## Developing

Once you've created a project and installed dependencies with `npm install` (or `pnpm install` or `yarn`), start a development server:

```bash
npm run dev

# or start the server and open the app in a new browser tab
npm run dev -- --open
```

## Building

To create a production version of your app:

```bash
npm run build
```

You can preview the production build with `npm run preview`.

> To deploy your app, you may need to install an [adapter](https://kit.svelte.dev/docs/adapters) for your target environment.


## Configuration 

Keycloak must match server side, please modify accordingly :  
- static/playground.html
- src/lib/keycloak.ts 

```
    keycloak = new Keycloak({
        url: 'https://your.url',
        realm: 'your_realm',
        clientId: 'your_client'
    });
```

## Note 

This line has been added : 
```
    checkLoginIframe: false
```
in 
```
    await keycloak.init({
        onLoad: 'login-required',
        checkLoginIframe: false
    });
```

else frontend will reload every 5 seconds
