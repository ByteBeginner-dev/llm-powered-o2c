/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
        mono: ['JetBrains Mono', 'monospace'],
      },
      colors: {
        'bg-base': 'var(--bg-base)',
        'bg-surface': 'var(--bg-surface)',
        'bg-elevated': 'var(--bg-elevated)',
        'border': 'var(--border)',
        'text-primary': 'var(--text-primary)',
        'text-muted': 'var(--text-muted)',
        'text-faint': 'var(--text-faint)',
        'accent': 'var(--accent)',
        'success': 'var(--success)',
        'warning': 'var(--warning)',
        'error': 'var(--error)',
        'node-order': 'var(--node-order)',
        'node-delivery': 'var(--node-delivery)',
        'node-billing': 'var(--node-billing)',
        'node-payment': 'var(--node-payment)',
        'node-customer': 'var(--node-customer)',
        'node-product': 'var(--node-product)',
        'node-plant': 'var(--node-plant)',
      },
      spacing: {
        '128': '32rem',
      },
      borderRadius: {
        'card': '8px',
      },
    },
  },
  plugins: [],
}
