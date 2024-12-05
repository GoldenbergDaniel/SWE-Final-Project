<script>
  import Authentication from "../lib/Authentication.svelte";
  import Background from "../lib/Background.svelte";
  import Footer from "../lib/Footer.svelte";

  let showSignIn = false;
  let showSignUp = false;

  const handleUpdateSignIn = (event) => {
    showSignIn = event.detail;
  };

  const handleUpdateSignUp = (event) => {
    showSignUp = event.detail;
  };
</script>

<main>
  <Background />
  <div class="content">
    <img src="/tradex_logo.jpg" class="logo" alt="TradEx Logo" />
    <h1>TradEx</h1>
    <p>Welcome to the Community Trading App</p>
  </div>
  
  <Authentication
    {showSignIn}
    {showSignUp}
    on:updateSignIn={handleUpdateSignIn}
    on:updateSignUp={handleUpdateSignUp}
  />

  <Footer />
</main>

<style>
  main {
    min-height: 550px;
  }

  main h1 {
    font-size: 4.5em;
  }

  main p {
    font-size: 1.5em;
  }

  .content {
    position: relative;
    z-index: 1;
    text-align: center;
  }

  .logo {
    height: 10em;
    transition: filter 300ms;
    will-change: filter;
  }

  .logo:hover {
    cursor: pointer;
  }
</style>