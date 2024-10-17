<script>
  const logoImage = './public/tradex_logo.jpg';
  import Authentication from "./lib/Authentication.svelte";

  let showMainPage = true;
  const handleUserData = (userData) => {
    showMainPage = false;
  };

</script>

<main>
  {#if showMainPage}
  <div>
    <img src={logoImage} class="logo" alt="TradEx Logo" />
  </div>
  <h1>TradEx</h1>
  <p>Welcome to the Community Trading App</p>

  <div class="user-info">
    <Authentication on:userData={handleUserData}/>
  </div>
  {:else if !showMainPage}
  <h1>NEW PAGE</h1>
  {/if}
</main>

<style>
  .logo {
    height: 10em;
    will-change: filter;
    transition: filter 300ms;
  }
  .logo:hover {
    filter: drop-shadow(0 0 5em #64ff6caa);
  }
</style>
