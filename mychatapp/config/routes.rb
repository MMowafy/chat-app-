Rails.application.routes.draw do
  # For details on the DSL available within this file, see https://guides.rubyonrails.org/routing.html
  resources :apps

  resources :apps do
    resources :chats
  end

  resources :apps do
    resources :chats do
      resources :messages
    end
  end

end
