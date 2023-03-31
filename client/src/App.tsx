import "./components/commonClasses.css"
import "./constants.css"
import Login from "./components/login/login"
import { Registration, AdditionalInfo } from "./components/registration/registration"
import { Route, Routes } from "react-router-dom"
import Navigation from "./components/header/header"
import PostList from "./components/user-post/postList"
import { UserInfo } from "./components/user-information/user-information"

// Shows the registration/login view, plus two buttons that lets you change between them
export default function App() {
    return (
        <Routes>
            <Route path="/login" element={<Login />} />
            <Route path="/registration" element={<Registration />} />
            <Route path="/additional-registration" element={<AdditionalInfo />} />
            <Route element={<Navigation />}>
                <Route path="/main" element={<h1>This is a main page</h1>} />
            </Route>
            <Route path="/user/:paramId" element={<UserInfo />} />
            <Route path="/user/:paramId/posts" element={<PostList />} />
            <Route path="/internal-error" element={<h1>Error 500</h1>} />
            <Route path="/*" element={<h1>Error 404</h1>} />
        </Routes>
    )
}
