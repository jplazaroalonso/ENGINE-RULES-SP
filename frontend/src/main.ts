import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { Quasar } from 'quasar'
import router from './router'

// Import icon libraries
import '@quasar/extras/material-icons/material-icons.css'
import '@quasar/extras/fontawesome-v6/fontawesome-v6.css'

// Import Quasar css
import 'quasar/src/css/index.sass'

// Import app styles
import './styles/main.scss'

// Import app component
import App from './App.vue'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(Quasar, {
  plugins: {}, // import Quasar plugins and add here
  config: {
    brand: {
      primary: '#004B87',    // Carrefour Blue
      secondary: '#E30613',  // Carrefour Red
      accent: '#FF6B35',     // Carrefour Orange
      positive: '#00A651',   // Carrefour Green
      negative: '#E30613',   // Carrefour Red
      info: '#004B87',       // Carrefour Blue
      warning: '#FF6B35'     // Carrefour Orange
    }
  }
})

app.mount('#app')
