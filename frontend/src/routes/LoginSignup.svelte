<script lang="ts">
  import * as Tabs from "$lib/components/ui/tabs";
  import Login from "$lib/components/login/Login.svelte";
  import Signup from "$lib/components/login/Signup.svelte";
</script>

<div class="container-sm mx-auto mt-6 sm:mt-12">
  <Tabs.Root value="login" class="sm:w-[400px] w-[350px]">
    <Tabs.List class="grid w-full grid-cols-2">
      <Tabs.Trigger value="login">Log In</Tabs.Trigger>
      <Tabs.Trigger value="signup">Sign Up</Tabs.Trigger>
    </Tabs.List>

    <Tabs.Content value="login">
      <Login />
    </Tabs.Content>
    <Tabs.Content value="signup">
      <Signup />
    </Tabs.Content>
  </Tabs.Root>
</div>
