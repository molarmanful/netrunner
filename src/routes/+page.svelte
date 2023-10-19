<script>
  import Cmd from '$lib/components/Cmd.svelte'
  import Env from '$lib/lang'

  let env = new Env()
</script>

<svelte:head>
  <title>netrunner</title>
</svelte:head>

{#each env.cmds as cmd}
  <Cmd {cmd} {env} />
{/each}
