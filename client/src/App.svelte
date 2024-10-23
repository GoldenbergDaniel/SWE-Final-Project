<script>
    import { Router, Link, Route } from 'svelte-routing';
    import Login from './pages/Login.svelte';
    import Dashboard from './pages/Dashboard.svelte';
    import Portfolio from './pages/Portfolio.svelte';
    import Connections from './pages/Connections.svelte';
    import Posts from './pages/Posts.svelte';
    import NotFound from './pages/404.svelte';
  </script>

  <Router>
    <Route path="/" component={Login} />
    <Route path="/dashboard" component={Dashboard} />
    <Route path="/portfolio" component={Portfolio} />
    <Route path="/connections" component={Connections} />
    <Route path="/posts" component={Posts} />
    <Route path="*" component={NotFound} />
  </Router>