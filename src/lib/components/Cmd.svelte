<script>
  import Btn from './Btn.svelte'

  export let env, k, f
</script>

<Btn
  on:click={() => {
    env.step(new env.types.Cmd(k), f)
  }}
>
  {k}
</Btn>
