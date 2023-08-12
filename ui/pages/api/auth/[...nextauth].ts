import NextAuth from "next-auth/next";
import CredentialsProvider from "next-auth/providers/credentials";

const backendURL = process.env.NEXT_PUBLIC_BACKEND_URL;
export const authOptions = {
session: {
    strategy: "jwt",
    maxAge: 30 * 24 * 60 * 60, // 30 days
  },
  providers: [
    CredentialsProvider({
      type: "credentials",
      credentials: {
        user: {
          label: "User",
          type: "text",
        },
        password: { label: "Password", type: "password" },
      },
      async authorize(credentials) {
        const credentialDetails = {
          username: credentials.user,
          password: credentials.password,
        };

        const resp = await fetch(backendURL + "/authentication/login", {
          method: "POST",
          headers: {
            Accept: "application/json",
            "Content-Type": "application/json",
          },
          body: JSON.stringify(credentialDetails),
        });
        const user = await resp.json();
        if (user.success) {
          console.log("nextauth daki user: " + user.success);

          return user;
        } else {
          console.log("check your credentials");
          console.log(user)
          return null;
        }
      },
    }),
  ],
callbacks: {
    jwt: async ({ token, user }) => {
      if (user) {
        //to do add groups
        token.username = user.response.name;
        token.accessToken = user.response.token[0];
      }

      return token;
    },
    session: ({ session, token, user }) => {
      if (token) {
        session.user.username = token.userName;
        session.user.accessToken = token.accessToken;
      }
      return session;
    },
  },
};

export default NextAuth(authOptions);