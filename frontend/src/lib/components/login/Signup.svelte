<script lang="ts">
  import * as Card from "$lib/components/ui/card";
  import Label from "../ui/label/label.svelte";
  import Input from "../ui/input/input.svelte";
  import FormItem from "./FormItem.svelte";
  import Button from "../ui/button/button.svelte";
  import { useSignup } from "./mutations";

  let email = "";
  let password = "";
  let password_confirmation = "";
  let name = "";

  const signup = useSignup();
</script>

<Card.Root>
  <Card.Header>
    <Card.Title>Sign Up</Card.Title>
    <Card.Description>Only SPA Emails will be accepted.</Card.Description>
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
        <Label>Name:</Label>
        <Input
          type="text"
          id="name"
          placeholder="Joe Spartan"
          bind:value={name}
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
      <FormItem>
        <Label>Password Confirmation:</Label>
        <Input
          type="password"
          id="password-confirm"
          placeholder="123445"
          bind:value={password_confirmation}
        />
      </FormItem>

      <Button
        class="mt-4"
        on:click={() =>
          $signup.mutate({
            email: email,
            name: name,
            password: password,
            password_confirm: password_confirmation,
          })}>Sign Up</Button
      >
    </form>
  </Card.Content>
</Card.Root>
