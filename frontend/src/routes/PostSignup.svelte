<script>
  import * as Card from "$lib/components/ui/card";
  import { Check } from "lucide-svelte";
</script>

<div class="container-sm mx-auto mt-6 sm:mt-12">
  <Card.Root class="w-[450px]">
    <Card.Header>
      <Card.Title>Welcome!</Card.Title>
      <Card.Description>This App Was Created by your SPA STC</Card.Description>
    </Card.Header>

    <!-- <Check class="text-green-400 w-20 h-20 mx-auto" /> -->

    <p class="text-left text-muted-foreground p-6 pt-0">
      You have recieved a verifcation email from our mailer, which is required
      to login to the application. If you do not recieve it, please check your
      spam folder, and contact the STC.
    </p>
  </Card.Root>
</div>
