<script>
  import { Btn, Cmd } from '$lib/components'
  import Env from '$lib/lang'

  let env = new Env()
  let { stack } = env

  let num = false
  let n = 0
</script>

<svelte:head>
  <title>netrunner</title>
</svelte:head>

<Btn
  on:click={() => {
    if (!num) num = true
    else {
      env.push(n)
      num = false
    }
  }}
>
  NUM
</Btn>

{#if num}
  {n}
  <input max="27" min="0" type="range" bind:value={n} />
{/if}

{#each Object.entries(env.cmds) as [k, f]}
  <Cmd {f} {k} />
{/each}

<br />
<br />
STACK:
<br />
<pre>
{$stack.join`\n`}
</pre>
