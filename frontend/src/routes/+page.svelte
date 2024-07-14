<script lang="ts">
  import store from '$lib/store';
  import type { PointageTime, Month } from '$lib/store';
  import { Login } from "../lib/wailsjs/go/main/App";
  import { goto } from '$app/navigation';
  import { slide } from 'svelte/transition';

  let email = '';
  let password = '';
  let loginError = false;
  let errorMessage: string;

  async function handleLogin() {
    const response = await Login(email, password);
    if (!response.success) {
      loginError = true;
      errorMessage = response.message;
      return;
    }

    store.set({ loginUser: response, pointages: new Map<Month, PointageTime[]>() });
    goto("/account");
  }
</script>

<div class="flex items-center justify-center min-h-screen bg-gray-100">
    <div class="w-full max-w-md p-6 bg-white rounded-lg shadow-md">
      <div class="space-y-2 text-center">
        <h1 class="font-semibold text-2xl">Connexion</h1>
        <p class="text-gray-500 dark:text-gray-400 text-sm">Entrez votre email et votre mot de passe pour accéder à votre compte</p>
      </div>
      <form on:submit|preventDefault={handleLogin}>
        <div class="form-group mb-1 flex flex-col">
          <label for="email" class=" text-zinc-600">Email</label>
          <input type="email" class="form-control border-[1px] rounded-md p-1 border-zinc-200 text-sm w-full py-2 focus:outline-none duration-300 transition-colors focus:bg-zinc-300 focus:border-zinc-400 { loginError ? " border-red-500" : "" }" id="email" placeholder="Email" bind:value={email} />
        </div>
        <div class="form-group mb-3 flex flex-col">
          <label for="password" class=" text-zinc-600">Password</label>
          <input type="password" class="form-control border-[1px] rounded-md p-1 border-zinc-200 text-sm w-full py-2 focus:outline-none duration-300 transition-colors focus:bg-zinc-300 focus:border-zinc-400 { loginError ? " border-red-500" : "" }" id="password" placeholder="Mot de passe" bind:value={password} />
        </div>
        {#if loginError}
          <div class="absolute right-2 bottom-2 bg-black w-60 text-sm p-5 shadow-lg rounded-lg text-zinc-200" role="alert" transition:slide>
            {errorMessage}
          </div>
        {/if}
        <button type="submit" class="btn btn-primary text-center bg-black rounded-md w-full h-full text-white p-2 duration-300 transition-colors  hover:bg-gray-800">Connexion</button>
      </form>
    </div>
</div>