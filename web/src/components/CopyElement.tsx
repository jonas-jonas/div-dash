import { IconClipboard } from "@tabler/icons";
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
      className="focus:outline-none hover:text-gray-900 inline-flex"
      onClick={handleClick}
    >
      <IconClipboard className="mr-1" size={16} />
      {value}
    </button>
  );
}
