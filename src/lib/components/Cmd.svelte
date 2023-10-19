<script>
  import Btn from './Btn.svelte'

  export let k, f
</script>

<Btn
  on:click={() => {
    f()
  }}
>
  {k}
</Btn>
