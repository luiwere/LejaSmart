function getColorFromKey(key) {
  const code = key.charCodeAt(0);
  return `hsl(${(code * 45) % 360}, 75%, 60%)`;
}

export function compose() {
  document.addEventListener('keydown', (event) => {
    const key = event.key;

    if (key >= 'a' && key <= 'z' && key.length === 1) {
      const div = document.createElement('div');
      div.className = 'note';
      div.textContent = key;
      div.style.backgroundColor = getColorFromKey(key);
      document.body.appendChild(div);
    } else if (key === 'Backspace') {
      const notes = document.querySelectorAll('.note');
      if (notes.length > 0) {
        notes[notes.length - 1].remove();
      }
    } else if (key === 'Escape') {
      document.querySelectorAll('.note').forEach(note => note.remove());
    }
  });
}
