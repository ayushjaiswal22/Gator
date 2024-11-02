package config

import(
    "encoding/json"
    "os"
)

const ConfigFile = ".gatorconfig.json"

type Config struct{
    DbUrl string `json:"db_url"`
    CurrentUsername string `json:"current_user_name"`
}

func Read() (Config, error) {

    home, err := os.UserHomeDir()
    if err!=nil {
        return Config{}, err
    }

    filePath := home + "/" +  ConfigFile
    if _, err := os.Stat(filePath); err!=nil {
        return Config{}, err
    }

    content, err := os.ReadFile(filePath)
    if err!=nil {
        return Config{}, err
    }
    var resp Config
    er := json.Unmarshal(content, &resp)
    if er!=nil {
        return Config{}, err
    }
    return resp, nil
}

func SetUser(conf Config, username string) error {
    conf.CurrentUsername = username
    jsonData, err := json.Marshal(conf)
    home, err := os.UserHomeDir()
    if err!=nil {
        return err
    }
    filePath := home + "/" +  ConfigFile
    os.WriteFile(filePath, jsonData, 0644)
    return nil
}
