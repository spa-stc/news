<script lang="ts">
  import * as Card from "$lib/components/ui/card";
  import Label from "../ui/label/label.svelte";
  import Input from "../ui/input/input.svelte";
  import FormItem from "./FormItem.svelte";
  import Button from "../ui/button/button.svelte";
  import { useLogin } from "./mutations";

  let email = "";
  let password = "";

  const login = useLogin();
</script>

<Card.Root>
  <Card.Header>
    <Card.Title>Log In</Card.Title>
    <Card.Description>Log in to your account here.</Card.Description>
  </Card.Header>
  <Card.Content class="space-y-2">
    <form on:submit|preventDefault>
      <FormItem>
        <Label>Email:</Label>
        <Input
          type="text"
          id="email"
          placeholder="stc@students.spa.edu"
          bind:value={email}
        />
      </FormItem>
      <FormItem>
        <Label>Password:</Label>
        <Input
          type="password"
          id="password"
          placeholder="123445"
          bind:value={password}
        />
      </FormItem>

      <Button
        class="mt-4"
        on:click={() => $login.mutate({ email: email, password: password })}
        >Log In</Button
      >
    </form>
  </Card.Content>
</Card.Root>
