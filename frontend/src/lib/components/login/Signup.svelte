<script lang="ts">
  import * as Card from "$lib/components/ui/card";
  import Label from "../ui/label/label.svelte";
  import Input from "../ui/input/input.svelte";
  import FormItem from "./FormItem.svelte";
  import Button from "../ui/button/button.svelte";
  import { pb } from "$lib/pocketbase";
  import { push } from "svelte-spa-router";
  import toast from "svelte-french-toast";
  import { ClientResponseError } from "pocketbase";

  let email = "";
  let password = "";
  let password_confirmation = "";
  let name = "";

  async function signup() {
    if (password != password_confirmation) {
      toast.error("passwords much match");
    }

    try {
      const data = {
        email: email,
        emailVisibility: false,
        password: password,
        passwordConfirm: password,
        name: name,
        role: "student",
      };

      const result = await pb.collection("users").create(data);

      await pb.collection("users").requestVerification(data.email);

      push("/");
    } catch (err) {
      if (err instanceof ClientResponseError) {
        toast.error(err.message);
        return;
      }

      toast.error("Something Went Wrong, Please Try Again.");
    }
  }
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

      <Button class="mt-4" on:click={signup}>Sign Up</Button>
    </form>
  </Card.Content>
</Card.Root>
