<script lang="ts">
  import { user } from "$lib/pocketbase";
  import * as Avatar from "$lib/components/ui/avatar";
  import Button from "../ui/button/button.svelte";
  import { push } from "svelte-spa-router";
</script>

{#if $user}
  <Avatar.Root>
    <Avatar.Image
      src="https://api.dicebear.com/8.x/thumbs/svg?seed={$user.name}"
    ></Avatar.Image>
    <Avatar.Fallback>{$user.name}</Avatar.Fallback>
  </Avatar.Root>
{:else}
  <Button on:click={() => push("/login")}>Login</Button>
{/if}
