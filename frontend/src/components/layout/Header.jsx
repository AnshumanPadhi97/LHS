import { Link } from 'react-router-dom';
import { Settings } from 'lucide-react';
import { Button } from "@heroui/button";

const Header = () => {
  return (
    <header className="bg-white border-b border-gray-200 shadow-sm py-4 px-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center">
          <Link to="/" className="text-xl font-bold text-blue-600">
            LHS - Local Hosting Services
          </Link>
        </div>
        <div className="flex items-center space-x-4">
          <Link to="/settings">
            <Button variant="outline" size="icon">
              <Settings className="h-5 w-5" />
            </Button>
          </Link>
        </div>
      </div>
    </header>
  );
};

export default Header;