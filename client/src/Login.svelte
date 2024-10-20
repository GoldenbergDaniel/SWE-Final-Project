<script>
  import Authentication from "./lib/Authentication.svelte";
  import Background from "./lib/Background.svelte";

  let showSignIn = false;
  let showSignUp = false;

  // Handlers for sign-in and sign-up events
  const handleUpdateSignIn = (event) => {
    showSignIn = event.detail;
  };

  const handleUpdateSignUp = (event) => {
    showSignUp = event.detail;
  };
</script>

<main>
  {#if !showSignIn && !showSignUp}
    <!-- Only show the background when neither SignIn nor SignUp are active -->
    <Background />
  {/if}
  <div class="content">
    <img src="/tradex_logo.jpg" class="logo" alt="TradEx Logo" />
    <h1>TradEx</h1>
    <p>Welcome to the Community Trading App</p>
  </div>
  
  <!-- Authentication Component -->
  <Authentication
    {showSignIn}
    {showSignUp}
    on:updateSignIn={handleUpdateSignIn}
    on:updateSignUp={handleUpdateSignUp}
  />
</main>

<style>
  .content {
    position: relative;
    z-index: 1; /* Make sure this content is above the background */
    text-align: center; /* Center the content */
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