<script>
  import { onMount } from 'svelte';

  let lastY = 0;

  const handleMouseMove = (e) => {
    const currentY = e.clientY;

    const dot = document.createElement('div');
    dot.classList.add('trail-dot');

    if (currentY < lastY) {
      dot.style.backgroundColor = '#4CAF50';
    } else {
      dot.style.backgroundColor = '#c3112c';
    }

    lastY = currentY;

    dot.style.left = `${e.pageX - 5}px`;
    dot.style.top = `${e.pageY - 5}px`;

    document.body.appendChild(dot);

    setTimeout(() => {
      dot.style.transform = 'scale(1.5)';
      dot.style.opacity = '0.5';
    }, 50);

    setTimeout(() => {
      dot.remove();
    }, 200);
  };

  // Initialize the event listener when the component mounts
  onMount(() => {
    document.addEventListener('mousemove', handleMouseMove);
    
    // Clean up the event listener when component unmounts
    return () => {
      document.removeEventListener('mousemove', handleMouseMove);
    };
  });
</script>