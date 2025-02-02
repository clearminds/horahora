import classNames from "classnames";
import { Link } from "react-router-dom";
import { Rate } from "antd";

const VIDEO_ELEMENT_WIDTH = "w-44";

function Video(props) {
  const { video } = props;

  return (
    <Link to={`/videos/${video.VideoID}`}>
      <div className={classNames(VIDEO_ELEMENT_WIDTH, "px-2 h-44")}>
        <div className="rounded overflow-hidden">
          <img
            className="block w-44 h-24 object-cover object-center"
            alt={video.Title}
            src={`${video.ThumbnailLoc}`}
            onError={(e)=>{e.target.onerror = null; e.target.src=`${video.ThumbnailLoc.slice(0, -6)}.jpg`}}
          />
          <Rate className="relative -mt-8 z-30 float-right" allowHalf={true} disabled={true} value={video.Rating}></Rate>
          </div>
        {/* TODO(ivan): deal with text truncation (hoping to have a multi-line text truncation,
                        which can't be done purely in css) */}
        <div className="text-xs font-bold text-blue-500 w-full py-1 text-black">{video.Title}</div>
        <div className="text-xs text-black">Views: {video.Views}</div>
      </div>
    </Link>
  );
}

function VideoList(props) {
  const { videos, title } = props;

  let elements = [];
  if (videos) {
    elements = [
      videos.map((video, idx) => <Video key={idx} video={video}/>),
    ];
  }

  // add padding elements so the items in the last row on this flexbox grid
  // get aligned with the other rows
  for (let i = 0; i < 10; i++) {
    elements.push(
      <div key={`_padding_${i}`} className={VIDEO_ELEMENT_WIDTH} />
    );
  }

  return (
    <div className="my-4 rounded border p-4 bg-white flex flex-wrap justify-around min-h-screen">
      {title ? <b className="text-xl">{title}</b> : <></>}
      {elements}
    </div>
  );
}

export default VideoList;
