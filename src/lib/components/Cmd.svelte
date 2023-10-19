<script>
  export let env
  export let cmd
</script>

<button on:click={env[cmd]()}>{cmd}</button>
