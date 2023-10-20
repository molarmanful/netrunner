<script>
  import { Btn, Cmd } from '$lib/components'
  import Env from '$lib/lang'

  let env = new Env()
  let { stack, code, scope, err } = env

  let num = false
  let get = false
  let set = false
  let n = 0
  let vname = ''
  let reset = () => {
    num = false
    get = false
    set = false
    n = 0
    vname = ''
  }
</script>

<svelte:head>
  <title>netrunner</title>
</svelte:head>

<main class="p-4">
  <div class="flex gap-4">
    <Btn
      on:click={() => {
        env.undo()
      }}
    >
      ←
    </Btn>

    <Btn
      on:click={() => {
        env.redo()
      }}
    >
      →
    </Btn>

    <Btn
      on:click={() => {
        if (!num) {
          reset()
          num = true
        } else reset()
      }}
    >
      NUM
    </Btn>

    <Btn
      on:click={() => {
        if (!get) {
          reset()
          get = true
        } else reset()
      }}
    >
      $
    </Btn>

    <Btn
      on:click={() => {
        if (!set) {
          reset()
          set = true
        } else reset()
      }}
    >
      $=
    </Btn>

    {#each Object.entries(env.cmds) as [k, f]}
      <Cmd {env} {f} {k} />
    {/each}
  </div>

  <br />

  {#if num}
    {n}
    <br />
    <input max="27" min="0" type="range" bind:value={n} />
    <Btn
      on:click={() => {
        env.step(n, () => {
          env.push(n)
        })
      }}
    >
      PUSH
    </Btn>
  {/if}

  {#if get}
    <input autofocus type="text" bind:value={vname} />
    <Btn
      on:click={() => {
        env.step('$' + vname, () => {
          let v = env.scope$[vname]
          if (!vname) throw new Error(`empty var`)
          if (v == void 0) throw new Error(`undefined var ${vname}`)
          env.push(v)
        })
        vname = ''
      }}
    >
      GET
    </Btn>
  {/if}

  {#if set}
    <input autofocus type="text" bind:value={vname} />
    <Btn
      on:click={() => {
        env.step('=$' + vname, () => {
          if (!vname) throw new Error(`empty var`)
          env.setVar(vname, env.pop())
        })
        vname = ''
      }}
    >
      SET
    </Btn>
  {/if}

  <hr />
  <br />
  {#if $err}
    <span class="c-red">ERR: {$err}</span>
  {/if}
  <pre>{($code, env.showCode())}</pre>
  <br />
  <div class="flex">
    <pre class="flex-1">{($stack, env.showStack())}</pre>
    <pre class="flex-1">{($scope, env.showScope())}</pre>
  </div>
</main>
