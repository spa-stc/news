<script lang="ts">
  import { user } from "$lib/pocketbase";
  import Button from "../ui/button/button.svelte";
  import { push } from "svelte-spa-router";
  import ProfileMenu from "./ProfileMenu.svelte";
</script>

{#if $user}
  <ProfileMenu username={$user.name} />
{:else}
  <Button on:click={() => push("/login")}>Login</Button>
{/if}
