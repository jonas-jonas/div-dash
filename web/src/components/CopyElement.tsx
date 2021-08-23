import {
  faCopy
} from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
type CopyElementProps = {
  value: string;
};

export function CopyElement({ value }: CopyElementProps) {
  const handleClick = () => {
    const textField = document.createElement("textarea");
    textField.innerText = value;
    document.body.appendChild(textField);
    textField.select();
    document.execCommand("copy");
    textField.remove();
  };

  return (
    <button
      className="focus:outline-none hover:text-gray-900"
      onClick={handleClick}
    >
      <FontAwesomeIcon icon={faCopy} className="mr-1" />
      {value}
    </button>
  );
}
