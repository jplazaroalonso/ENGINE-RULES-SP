# Rules Engine Frontend

A modern Vue 3 Single Page Application (SPA) for managing business rules with Carrefour-inspired design.

## 🎨 Design System

This application uses a comprehensive design system inspired by Carrefour's brand identity:

- **Primary Color**: Carrefour Blue (#004B87)
- **Secondary Color**: Carrefour Red (#E30613)
- **Accent Color**: Carrefour Orange (#FF6B35)
- **Success Color**: Carrefour Green (#00A651)
- **Typography**: Inter font family
- **Icons**: Material Design Icons

## 🚀 Features

### Core Functionality
- **Rules Management**: Create, edit, delete, and manage business rules
- **Rule Evaluation**: Test rules with sample data
- **Rule Calculation**: Calculate results from multiple rules
- **Campaign Management**: Manage promotional campaigns
- **Analytics Dashboard**: View performance metrics and insights
- **Customer Management**: Manage customer data and segments

### User Experience
- **Responsive Design**: Mobile-first approach with progressive enhancement
- **Dark/Light Mode**: Toggle between themes
- **Real-time Updates**: Live data synchronization
- **Offline Support**: Progressive Web App capabilities
- **Accessibility**: WCAG 2.1 AA compliance

### Technical Features
- **TypeScript**: Full type safety throughout the application
- **Component Library**: Reusable UI components
- **State Management**: Pinia for reactive state management
- **API Integration**: RESTful API client with error handling
- **Testing**: Unit, integration, and E2E tests
- **Performance**: Optimized bundle size and loading times

## 🛠️ Technology Stack

- **Framework**: Vue 3 with Composition API
- **Language**: TypeScript
- **UI Library**: Quasar Framework
- **State Management**: Pinia
- **Routing**: Vue Router
- **HTTP Client**: Axios
- **Charts**: Chart.js with Vue-ChartJS
- **Code Editor**: Monaco Editor
- **Build Tool**: Vite
- **Testing**: Vitest, Cypress
- **Linting**: ESLint, Prettier
- **Styling**: SCSS with CSS Custom Properties

## 📦 Installation

### Prerequisites
- Node.js 18.0.0 or higher
- npm or yarn package manager

### Setup
1. Clone the repository
2. Navigate to the frontend directory
3. Install dependencies:
   ```bash
   npm install
   ```

4. Copy environment variables:
   ```bash
   cp env.example .env
   ```

5. Update the `.env` file with your configuration

### Development
```bash
# Start development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview

# Run tests
npm run test

# Run E2E tests
npm run test:e2e

# Lint code
npm run lint

# Format code
npm run format
```

## 🏗️ Project Structure

```
frontend/
├── public/                 # Static assets
│   ├── favicon.svg        # Carrefour-inspired favicon
│   └── index.html         # HTML template
├── src/
│   ├── api/               # API client and services
│   │   ├── client.ts      # Axios configuration
│   │   └── rules.ts       # Rules API methods
│   ├── components/        # Vue components
│   │   ├── common/        # Reusable components
│   │   ├── charts/        # Chart components
│   │   └── layout/        # Layout components
│   ├── router/            # Vue Router configuration
│   ├── stores/            # Pinia stores
│   │   ├── auth.ts        # Authentication store
│   │   ├── rules.ts       # Rules management store
│   │   ├── notifications.ts # Notifications store
│   │   └── ui.ts          # UI state store
│   ├── styles/            # SCSS styles
│   │   ├── main.scss      # Main stylesheet
│   │   └── quasar-variables.scss # Quasar overrides
│   ├── types/             # TypeScript type definitions
│   ├── utils/             # Utility functions
│   ├── views/             # Page components
│   │   ├── auth/          # Authentication pages
│   │   ├── rules/         # Rules management pages
│   │   └── Dashboard.vue  # Main dashboard
│   ├── App.vue            # Root component
│   └── main.ts            # Application entry point
├── package.json           # Dependencies and scripts
├── vite.config.ts         # Vite configuration
├── tsconfig.json          # TypeScript configuration
└── README.md              # This file
```

## 🎯 Key Components

### Layout Components
- **AppHeader**: Top navigation with search, notifications, and user menu
- **MainNavigation**: Sidebar navigation with role-based menu items
- **AppFooter**: Footer with system status and links

### Common Components
- **MetricCard**: Displays key performance indicators
- **PerformanceChart**: Interactive charts for analytics
- **RuleCard**: Displays rule information in a card format

### Views
- **Dashboard**: Main overview with metrics and recent activity
- **RulesList**: List and manage all business rules
- **RuleEditor**: Create and edit rules with DSL editor
- **Login**: Authentication page with Carrefour branding

## 🔧 Configuration

### Environment Variables
- `VITE_API_BASE_URL`: Base URL for API requests
- `VITE_API_TIMEOUT`: Request timeout in milliseconds
- `VITE_APP_NAME`: Application name
- `VITE_ENABLE_ANALYTICS`: Enable analytics tracking
- `VITE_ENABLE_DEBUG_MODE`: Enable debug logging

### API Integration
The application integrates with three main services:
- **Rules Management Service**: CRUD operations for rules
- **Rules Evaluation Service**: Rule testing and validation
- **Rules Calculator Service**: Rule calculation and execution

## 🎨 Design System

### Color Palette
```scss
// Primary colors
--carrefour-blue: #004B87;
--carrefour-red: #E30613;
--carrefour-orange: #FF6B35;
--carrefour-green: #00A651;

// Light variants
--carrefour-light-blue: #E3F2FD;
--carrefour-light-red: #FFEBEE;
--carrefour-light-orange: #FFF3E0;
--carrefour-light-green: #E8F5E8;

// Neutral grays
--carrefour-gray-50: #FAFAFA;
--carrefour-gray-100: #F5F5F5;
--carrefour-gray-200: #EEEEEE;
--carrefour-gray-300: #E0E0E0;
--carrefour-gray-400: #BDBDBD;
--carrefour-gray-500: #9E9E9E;
--carrefour-gray-600: #757575;
--carrefour-gray-700: #616161;
--carrefour-gray-800: #424242;
--carrefour-gray-900: #212121;
```

### Typography
- **Font Family**: Inter (primary), Courier New (monospace)
- **Font Sizes**: 12px to 36px scale
- **Font Weights**: 300, 400, 500, 600, 700
- **Line Heights**: 1.25, 1.5, 1.75

### Spacing
- **Scale**: 4px base unit (4px, 8px, 12px, 16px, 20px, 24px, 32px, 40px, 48px, 64px, 80px, 96px)

### Border Radius
- **Small**: 4px
- **Medium**: 8px
- **Large**: 12px
- **Extra Large**: 16px
- **Full**: 50%

## 🧪 Testing

### Unit Tests
```bash
npm run test
```

### E2E Tests
```bash
npm run test:e2e
```

### Test Coverage
- Components: 80%+ coverage
- Stores: 90%+ coverage
- Utils: 95%+ coverage

## 🚀 Deployment

### Production Build
```bash
npm run build
```

### Docker
```dockerfile
FROM node:18-alpine as build
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

### Environment Setup
1. Set production environment variables
2. Configure API endpoints
3. Enable analytics and monitoring
4. Set up CDN for static assets

## 📱 Progressive Web App

The application includes PWA features:
- **Service Worker**: Offline functionality
- **Web App Manifest**: App-like experience
- **Push Notifications**: Real-time updates
- **Background Sync**: Data synchronization

## ♿ Accessibility

- **WCAG 2.1 AA Compliance**: Meets accessibility standards
- **Keyboard Navigation**: Full keyboard support
- **Screen Reader Support**: ARIA labels and descriptions
- **High Contrast Mode**: Support for high contrast displays
- **Reduced Motion**: Respects user motion preferences

## 🔒 Security

- **HTTPS Only**: All communications encrypted
- **Content Security Policy**: Prevents XSS attacks
- **Authentication**: JWT-based authentication
- **Authorization**: Role-based access control
- **Input Validation**: Client and server-side validation

## 📊 Performance

### Metrics
- **First Contentful Paint**: < 1.5s
- **Largest Contentful Paint**: < 2.5s
- **Cumulative Layout Shift**: < 0.1
- **First Input Delay**: < 100ms

### Optimization
- **Code Splitting**: Lazy loading of routes
- **Tree Shaking**: Remove unused code
- **Image Optimization**: WebP format with fallbacks
- **Caching**: Aggressive caching strategy

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

For support and questions:
- **Documentation**: Check the project documentation
- **Issues**: Create an issue on GitHub
- **Email**: Contact the development team

## 🔄 Version History

- **v1.0.0**: Initial release with core functionality
- **v1.1.0**: Added analytics dashboard
- **v1.2.0**: Enhanced rule editor with syntax highlighting
- **v1.3.0**: Added campaign management features

---

Built with ❤️ for Carrefour using Vue 3 and modern web technologies.
