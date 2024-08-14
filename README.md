# QChat - Messaging Web Application
qChat is a scalable chat application built using React. This frontend application provides a user-friendly interface for real-time messaging, user authentication, and chat history.

## Features

- **Real-time Messaging:** Supports instant messaging between users with live updates.
- **User Authentication:** Secure login and registration process.
- **Chat History:** Ability to view previous chat history.
- **Offline Chat:** Facility to chat even users are offline
- **Attractive Design:** Adaptable UI design with Bootstrap design guidelines.

## Technologies Used

- **React:** The primary JavaScript library used for building the user interface.
- **TypeScript:** A superset of JavaScript used to add type definitions, improving code quality.
- **Vite:** A build tool used for bundling the application.
- **React Router:** Handles routing within the application.
- **Axios:** A promise-based HTTP client used for making API requests to the backend.
- **Bootstrap:** A CSS framework used to design responsive layouts.
- **WebSockets:** Facilitates real-time, bidirectional communication between the client and server.

## Project Setup

To set up the project locally, follow these steps:

1. **Clone the repository**
```bash
   git clone https://github.com/darkzChira/qchat.git
   cd qchat
```

2. **Install dependencies**
```bash
npm install
```

3. **Configure API**
   Add the Backend API URL to `vite.config.ts` file
```bash
proxy: {
      '/chatapp': {
        target: <Add BE URL>,
        changeOrigin: true,
        ws: true,
        rewrite: (path) => path.replace(/^\/chatapp/, ''),
      },
```


4. Run the application
```bash
npm run dev
```

The development server will start, typically accessible at http://localhost:5173/.

## Deployment
The application has been deployed to Vercel.

The production application is accessible at [https://qchatfe.vercel.app/](https://qchatfe.vercel.app/).

## Contact
For any inquiries or feedback, feel free to reach out:

* Email: [darakachiranjaya@gmail.com](darakachiranjaya@gmail.com)
* LinkedIn: [https://www.linkedin.com/in/daraka-chiranjaya/](https://www.linkedin.com/in/daraka-chiranjaya/)


