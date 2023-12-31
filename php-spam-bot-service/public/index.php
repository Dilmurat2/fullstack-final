<?php

require_once __DIR__ . '/../vendor/autoload.php';

use App\Config;
use App\Http\Controllers;
use Doctrine\ORM\EntityManager;
use Doctrine\ORM\ORMSetup;
use Dotenv\Dotenv;
use Slim\Factory\AppFactory;
use Slim\Views\Twig;
use Slim\Views\TwigMiddleware;
use Twig\Error\LoaderError;
use DI\Container;
use function DI\create;

try {
    $twig = Twig::create(__DIR__ . '/../template', ['cache' => false]);
} catch (LoaderError $e) {
    die(500);
}

$dotenv = Dotenv::createImmutable(dirname(__DIR__));
$dotenv->load();

//dependency injection
$container = new Container();
//doctrine
$container->set(Config::class, create(Config::class)->constructor($_ENV));
$container->set(EntityManager::class, fn(Config $config) => EntityManager::create(
    $config->db,
    ORMSetup::createAttributeMetadataConfiguration([__DIR__ . '/../src/Model'])
));

AppFactory::setContainer($container);

$app = AppFactory::create();

$app->addErrorMiddleware(true, true, true);

$app->add(TwigMiddleware::create($app, $twig));

$app->get('/history', Controllers\IndexController::class . ':history');
$app->get('/', Controllers\MessageController::class . ':createMessage');

$app->post('/', Controllers\MessageController::class . ':sendMessage');

$app->run();
