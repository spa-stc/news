<script lang="ts">
  import * as Avatar from "$lib/components/ui/avatar";
  import * as DropdownMenu from "$lib/components/ui/dropdown-menu";
  import { pb } from "$lib/pocketbase";
  import toast from "svelte-french-toast";

  export let username: string;
</script>

<DropdownMenu.Root>
  <DropdownMenu.Trigger>
    <Avatar.Root>
      <Avatar.Image
        src="https://api.dicebear.com/8.x/thumbs/svg?seed={username}"
      ></Avatar.Image>
      <Avatar.Fallback>{username}</Avatar.Fallback>
    </Avatar.Root>
  </DropdownMenu.Trigger>
  <DropdownMenu.Content>
    <DropdownMenu.Group>
      <DropdownMenu.Label>Account</DropdownMenu.Label>
      <DropdownMenu.Item
        on:click={() => {
          pb.authStore.clear();
          toast.success("Logged Out!");
        }}>Log Out</DropdownMenu.Item
      >
    </DropdownMenu.Group>
  </DropdownMenu.Content>
</DropdownMenu.Root>
