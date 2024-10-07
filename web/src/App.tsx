import './App.css'
import {useState} from "react";

export const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080/v1"

export default function App() {
    return (
        <div>
            <header className="header">
                <div className="container header-content">
                    <a href="/" className="logo">&lt;Gossage /&gt;</a>
                    <nav className="nav">
                        <a href="/" className="nav-link">Home</a>
                        <a href="/about" className="nav-link">About</a>
                        <a href="/contact" className="nav-link">Contact</a>
                    </nav>
                </div>
            </header>

            <main className="container main">
                <div className="content">
                    <article>
                        <h2>Sample Blog Post</h2>
                        <p>Published on June 1, 2023</p>
                        <p>
                            This is a sample blog post. Your actual blog content would go here.
                            You can add multiple paragraphs, code snippets, and other elements as needed.
                        </p>
                        <pre>
              <code>
{`function helloWorld() {
  console.log("Hello, World!");
}

helloWorld();`}
              </code>
            </pre>
                    </article>
                    <button className="button">Read More</button>
                </div>

            </main>

            <footer className="footer">
                <div className="container">
                    <p>&copy; 2023 MyBlog. All rights reserved.</p>
                </div>
            </footer>
        </div>
    );
}
