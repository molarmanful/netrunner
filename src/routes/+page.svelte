<script>
  import { Btn, Cmd } from '$lib/components'
  import Env from '$lib/lang'

  let env = new Env()
  let { stack, code, err } = env

  let num = false
  let n = 0
</script>

<svelte:head>
  <title>netrunner</title>
</svelte:head>

<main class="p-4">
  <div class="flex gap-4">
    <Btn
      on:click={() => {
        if (!num) num = true
        else {
          env.step(() => {
            env.push(n)
          })
          num = false
        }
      }}
    >
      NUM
    </Btn>
    {#each Object.entries(env.cmds) as [k, f]}
      <Cmd {env} {f} {k} />
    {/each}
  </div>

  <br />

  {#if num}
    {n}
    <input max="27" min="0" type="range" bind:value={n} />
  {/if}

  <hr />
  <br />
  {#if $err}
    ERR: {$err}
  {/if}
  <pre>
{($stack, env.showStack())}
  </pre>
</main>
