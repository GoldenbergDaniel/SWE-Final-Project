<script>
    import { Router, Link, Route } from 'svelte-routing';
    import Login from './pages/Login.svelte';
    import Portfolio from './pages/Portfolio.svelte';
    import Trade from './pages/Trade.svelte';
    import Posts from './pages/Posts.svelte';
    import Leaderboard from './pages/Leaderboard.svelte';
    import NotFound from './pages/404.svelte';
  </script>

  <Router>
    <Route path="/" component={Login} />
    <Route path="/portfolio" component={Portfolio} />
    <Route path="/trade" component={Trade} />
    <Route path="/posts" component={Posts} />
    <Route path="/leaderboard" component={Leaderboard} />
    <Route path="*" component={NotFound} />
  </Router>