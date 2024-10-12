<script>
  const logoImage = './public/tradex_logo.jpg';
  import Authentication from "./lib/Authentication.svelte";
</script>

<main>
  <div>
    <a href="https://vitejs.dev" target="_blank" rel="noreferrer">
      <img src={logoImage} class="logo" alt="TradEx Logo" />
    </a>
  </div>
  <h1>TradEx</h1>
  <p>Welcome to the Community Trading App</p>

  <div class="user-info">
    <Authentication />
  </div>

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
