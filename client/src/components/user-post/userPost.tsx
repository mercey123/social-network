import { useEffect, useState } from "react"
import icon from "../../assets/SVGRepo_iconCarrier.svg"
import { useOpenText } from "../../hooks/openText"
import { Post } from "../models"
import CommentList from "./comment"
import Loading from "./render-states/loading"
import "./userPost.css"
import { getPostData } from "./fetch"

export default function UserPost({ postId }: { postId: number }) {
    const [post, setPost] = useState<Post | undefined>(undefined)
    const [err, setErr] = useState<Error | null>(null)
    const [isLoading, setIsLoading] = useState(true)
    const [attachmentsCount, setAttachmentsCount] = useState(0)
    const [attachmentsOpen, setAttachmentsOpen] = useState(false)
    const [commentsOpen, setCommentsOpen] = useState(false)
    const { style, refText, openText } = useOpenText(215)

    useEffect(() => {
        getPostData(postId, setPost, setErr, setIsLoading)
    }, [])

    useEffect(() => {
        if (post) {
            setAttachmentsCount(post.attachments?.length)
        }
    }, [post])

    const toggleAttachments = () => setAttachmentsOpen(!attachmentsOpen)

    if (err) return <div>{err.message}</div>
    if (isLoading) return <Loading color="orange" />
    return (
        <article className="post">
            {post.text ? (
                <p
                    ref={refText}
                    className="post__text"
                    style={style}
                    onClick={() => {
                        openText()
                    }}
                >
                    {post.text}
                </p>
            ) : (
                <img className="post__attachment" src={post.attachments[0]} alt="placeholder" />
            )}

            {attachmentsCount > 0 &&
                attachmentsOpen &&
                post.attachments
                    .slice(post.text ? 0 : 1)
                    .map((att, i) => (
                        <img className="post__attachment" key={i} src={att} alt="placeholder" />
                    ))}

            <div className="post__actions">
                {attachmentsCount > 0 && (
                    <p className="post__attachments-toggler" onClick={() => toggleAttachments()}>
                        {attachmentsOpen
                            ? "Show less"
                            : `Show ${
                                  post.text ? attachmentsCount : attachmentsCount - 1
                              } more attachment${attachmentsCount > 1 ? "s" : ""}`}
                    </p>
                )}

                <img
                    className="post__add-comment"
                    src={icon}
                    alt="comment icon"
                    onClick={() => setCommentsOpen(!commentsOpen)}
                />
            </div>

            {commentsOpen && <CommentList postId={postId} />}
        </article>
    )
}
